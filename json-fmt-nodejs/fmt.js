const fs = require('fs');

try {
    const data = fs.readFileSync('Sofia-API-Meta-Data-18072023.json', 'utf8');
    const obj = JSON.parse(data)
    const str = JSON.stringify(obj, null, 4)
    fs.writeFileSync('Sofia-API-Meta-Data-18072023-fmt.json', str);
} catch (e) {
    console.error(e);
}