function markdownToHtmlTable(markdown, tableName) {
    const lines = markdown.trim().split('\n');
    if (lines.length < 2) return '';

    // Extract headers
    const headers = parseMarkdownLine(lines[0]);
    // Extract rows
    const rows = lines.slice(2).map(line => parseMarkdownLine(line));

    let html = `<table><caption>${tableName}</caption><thead><tr>`;

    // Create table header
    headers.forEach(header => {
        html += `<th>${header}</th>`;
    });
    html += '</tr></thead><tbody>';

    // Create table rows
    rows.forEach((row, index) => {
        html += '<tr>';
        row.forEach(cell => {
            html += `<td>${cell}</td>`;
        });
        html += '</tr>';
    });

    html += '</tbody></table>';
    return html;
}


function parseMarkdownLine(line) {
    // Regex para identificar e manter pipes escapados
    const regex = /\\\|/g;
    let parts = [];
    let part = '';
    let escape = false;

    // Itera sobre cada caractere da linha
    for (let char of line) {
        if (escape) {
            part += char;  // Adiciona o caractere escapado
            escape = false;
        } else if (char === '\\') {
            escape = true;  // Marca que o próximo caractere deve ser escapado
        } else if (char === '|') {
            parts.push(part.trim());  // Adiciona a parte até o pipe
            part = '';  // Reinicia a parte
        } else {
            part += char;  // Adiciona o caractere normal
        }
    }

    // Adiciona a última parte
    if (part.length > 0) {
        parts.push(part.trim());
    }

    return parts;
}

document.addEventListener('DOMContentLoaded', function() {
    // Obtém o nome da tabela a partir do atributo data do script
    var tableName = document.querySelector('script[data-table-name]').getAttribute('data-table-name');

    // Obtém o conteúdo da tabela do div
    var markdownContent = document.getElementById('table-container').textContent;

    // Converte Markdown para HTML
    var tableHtml = markdownToHtmlTable(markdownContent, tableName);

    // Substitui o conteúdo do div pelo HTML da tabela com caption
    document.getElementById('table-container').innerHTML = tableHtml;
});

