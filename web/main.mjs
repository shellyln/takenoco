import './wasm_exec.js';
import { readFileSync } from 'fs';

const go = new Go();
const mod = await WebAssembly.compile(readFileSync('./go.wasm'));
let inst = await WebAssembly.instantiate(mod, go.importObject);

// NOTE: TinyGo can export symbols.
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
