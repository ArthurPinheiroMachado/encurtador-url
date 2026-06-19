FROM mcr.microsoft.com/dotnet/sdk:8.0 AS build
WORKDIR /src
COPY Encurtador.Url/Encurtador.Url.csproj Encurtador.Url/
RUN dotnet restore Encurtador.Url/Encurtador.Url.csproj
COPY . .
WORKDIR /src/Encurtador.Url
RUN dotnet publish -c Release -o /app

FROM mcr.microsoft.com/dotnet/aspnet:8.0
WORKDIR /app
COPY --from=build /app .
EXPOSE 6060
ENTRYPOINT ["dotnet", "Encurtador.Url.dll"]
