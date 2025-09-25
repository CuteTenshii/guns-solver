const workerScript = `        
            import init, { GunsSolver } from 'https://assets.guns.lol/pkg/_gunslolpow.js';
            self.onmessage = async function (event) {               
                const { ps, d, c, _n, org_ts } = event.data;
                await init();
                const _get_s = new GunsSolver(ps, c, parseInt(d), false, _n);
                const ts = Math.floor(Date.now() / 1000);
                const _tsoff = ts - Number(org_ts);
                const tsp = ts - _tsoff;
                const _res = await _get_s.solve_pow();    
                if (_res !== null) {
                    self.postMessage({ _res, ps, tsp, _n, c });
                }
            };
        `;

const blob = new Blob([workerScript], { type: "application/javascript" });
const worker = new Worker(URL.createObjectURL(blob), { type: "module" });

worker.postMessage({
    ps: '1b9e851d85c31e0b036b1d822c9692b962a7851086758b740f3643c1376',
    d: '5',
    c: '1a1e20b1e8c66e81735d8b892abc2339415635871fd7da0d1908fdaab4c860a0',
    _n: 'eqkVruB41bP6oeCn',
    org_ts: '1758726615'
});

worker.onmessage = async function (event) {
    const { _res, ps, tsp, _n, c } = event.data;
    if (_res !== undefined) {
        await getResult(_res);
    }
};