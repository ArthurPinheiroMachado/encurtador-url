from contextlib import asynccontextmanager
from urllib.parse import urlparse

from fastapi import FastAPI, Depends, HTTPException
from fastapi.responses import JSONResponse, RedirectResponse, Response

from .auth import verify_auth
from .cache import url_cache
from .config import settings
from .database import get_all_urls, get_url_by_original, save_url, delete_url, update_accesses, migrate, close_pool
from .models import Info, UrlCreate, UrlCreated
from .utils import generate_short_id


@asynccontextmanager
async def lifespan(app: FastAPI):
    await migrate()
    urls = await get_all_urls()
    await url_cache.load(urls)
    yield
    await close_pool()


prefix = settings.HTTP_BASE.rstrip("/")
app = FastAPI(lifespan=lifespan)


@app.get(f"{prefix}/urls", dependencies=[Depends(verify_auth)])
async def get_urls():
    return await url_cache.get_all()


@app.post(f"{prefix}/urls", status_code=201, dependencies=[Depends(verify_auth)])
async def create_url(payload: UrlCreate):
    parsed = urlparse(payload.url)
    if not parsed.scheme or not parsed.netloc:
        raise HTTPException(status_code=400, detail="Invalid URL")

    existing = await get_url_by_original(payload.url)
    if existing:
        return JSONResponse(content={"id": existing["id"], "url": payload.url}, status_code=200)

    short_id = generate_short_id(8, lambda id: url_cache.id_exists(id))

    await save_url(short_id, payload.url)
    await url_cache.set(short_id, Info(original=payload.url, accesses=0))

    return UrlCreated(id=short_id, url=payload.url)


@app.get(f"{prefix}/urls/{{id}}", dependencies=[Depends(verify_auth)])
async def get_url(id: str):
    info = await url_cache.get(id)
    if info is None:
        raise HTTPException(status_code=400, detail="URL not found")
    return info


@app.get(f"{prefix}/{{id}}", dependencies=[Depends(verify_auth)])
async def redirect_to_original(id: str):
    info = await url_cache.get(id)
    if info is None:
        raise HTTPException(status_code=404, detail="URL not found")

    new_accesses = await url_cache.increment_accesses(id)
    await update_accesses(id, new_accesses)

    return RedirectResponse(url=info.original, status_code=302)


@app.delete(f"{prefix}/{{id}}", dependencies=[Depends(verify_auth)])
async def delete_url_route(id: str):
    exists = await url_cache.exists(id)
    if not exists:
        raise HTTPException(status_code=400, detail="URL not found")

    await delete_url(id)
    await url_cache.delete(id)
    return Response(status_code=200)
