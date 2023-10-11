const fs = require('fs');

// Read the arguments from the command line
const fPath = process.argv[2];

try {

    const data = fs.readFileSync(fPath, 'utf8');
    const str = JSON.stringify(JSON.parse(data), null, 2)
    // fs.writeFileSync('Sofia-API-Meta-Data-18072023-fmt.json', str);    
    console.log(str)

} catch (e) {
    console.error(e);
}

