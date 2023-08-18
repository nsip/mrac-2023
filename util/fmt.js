// Read the arguments from the command line
const strJSON = process.argv[2];

// Print the result
console.log(JSON.stringify(JSON.parse(strJSON), null, 2));
