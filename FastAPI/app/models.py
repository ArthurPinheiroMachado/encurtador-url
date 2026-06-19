from pydantic import BaseModel


class UrlRecord(BaseModel):
    id: str
    original: str
    accesses: int


class Info(BaseModel):
    original: str
    accesses: int


class UrlCreate(BaseModel):
    url: str


class UrlCreated(BaseModel):
    id: str
    url: str
