import { readFileSync } from 'fs';
import * as url from 'url';

globalThis.crypto = await import('node:crypto');
await import('./wasm_exec.js');

const __filename = url.fileURLToPath(import.meta.url);
const __dirname = url.fileURLToPath(new URL('.', import.meta.url));

const go = new Go();
const mod = await WebAssembly.compile(readFileSync(__dirname + '/go.wasm'));
let inst = await WebAssembly.instantiate(mod, go.importObject);

globalThis.goWasmExports = inst.exports;

async function run() {
    console.clear();

    const goInstRan = go.run(inst);

    {
        let result = parseCsv('0,1,2,3,4,5,6,7,8,9\n0,1,2,3,4,5,6,7,8,9');
        console.log(result);
    }

    await goInstRan;
}

await run();
