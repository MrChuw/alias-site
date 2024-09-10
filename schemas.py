from tortoise import fields
from tortoise.models import Model


class MarkdownTable(Model):
    id = fields.CharField(max_length=32, unique=True, primary_key=True)
    content = fields.TextField()
    table_name = fields.TextField(max_length=255)

    class Meta:
        table = "markdown_tables"
