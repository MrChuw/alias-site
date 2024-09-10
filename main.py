import uuid

from fastapi import FastAPI, Form, HTTPException, Request
from fastapi.responses import HTMLResponse, RedirectResponse
from fastapi.templating import Jinja2Templates
from tortoise.contrib.fastapi import register_tortoise
from fastapi.staticfiles import StaticFiles
from schemas import MarkdownTable
from custom_logging import CustomizeLogger
from pathlib import Path
import uvicorn
import logging



logger = logging.getLogger(__name__)
config_path = Path(__file__).with_name("logging_config.json")
app = FastAPI(debug=False)
app.logger = CustomizeLogger.make_logger(config_path)
templates = Jinja2Templates(directory="templates")
app.mount("/static", StaticFiles(directory="static"), name="static")


@app.get("/", response_class=HTMLResponse)
async def index(request: Request):
    return templates.TemplateResponse("index.html", {"request": request})


@app.get("/favicon.ico")
async def get_favicon():
    return app.url_path_for('static', path='favicon.ico')


@app.post("/submit")
async def submit_table(request: Request, markdown: str = Form(...), table_name: str = Form(...)):
    if not markdown or not table_name:
        raise HTTPException(status_code=400, detail="No Markdown content or table name provided")

    table_id = uuid.uuid4().hex[:16]
    await MarkdownTable.create(id=table_id, content=markdown, table_name=table_name)
    return RedirectResponse(url=f"/{table_id}", status_code=303)



@app.get("/{table_id}", response_class=HTMLResponse)
async def view_table(request: Request, table_id: str):
    table = await MarkdownTable.get_or_none(id=table_id)
    if not table:
        raise HTTPException(status_code=404, detail="Alias not found")

    return templates.TemplateResponse("table.html",
                                      {"request": request, "Content": table.content, "TableName": table.table_name}
                                      )



register_tortoise(app, db_url="sqlite://db.sqlite3", modules={"models": ["schemas"]}, generate_schemas=True,
                  add_exception_handlers=True)

if __name__ == "__main__":
    # uvicorn.run("main:app", host="0.0.0.0", port=5454, log_level="info", reload=True)
    uvicorn.run("main:app", host="0.0.0.0", port=5454, log_level="info", workers=4)
