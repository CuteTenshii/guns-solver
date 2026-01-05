// language=JavaScript
const workerScript = `
    import init, {GunsSolver} from 'https://assets.guns.lol/wasm/gpp_gunslol.js';

    self.onmessage = async function (event) {
        const {_2xa, _org_ts, _n, d, o09} = event.data;
        await init();
        const _get_s = new GunsSolver(o09, parseInt(d), _org_ts, _n, _2xa);
        const ts = Math.floor(Date.now() / 1000);
        const _tsoff = ts - Number(_org_ts);
        const tsp = ts - _tsoff;
        const _res = await _get_s.solve_pow();
        if (_res && _res._oo && _res.seal) {
            self.postMessage({_n, o09, _res, _2xa, tsp});
        }
    };`;
const blob = new Blob([workerScript], {type: "application/javascript"});
const worker = new Worker(URL.createObjectURL(blob), {type: "module"});
const _gs_sets = {
    _n: 'Zn6F5oHCaRnX1vT9zQ81a4Ty7faAAGCb',
    _tsp: null,
    _p: '',
    _seal: '',
    __ps: '',
    o09: '1867764c175fc766ce086157515853444b999d7f64d03cbe1d96b78cf6224e02',
    _2xa: 'oUAFFCgIIAUAAgEDBFRWE1Bbr_TxYzk5MjZmYWI4NTZkZjlhNjQ5YmY1MWM1NWUyOTk5ZDJmZjE1Y2ZiZGNjYjhjY2Q1YzgzNzMzNTc2NGIXqZPHa_eHIw',
    __s: '3',
    _org_ts: "1767638499",
};
worker.postMessage({
    _n: _gs_sets._n,
    d: 5,
    _2xa: _gs_sets._2xa,
    o09: _gs_sets.o09,
    _org_ts: _gs_sets._org_ts
});
worker.onmessage = async function (event) {
    const {_res, o09, tsp, _n, _2xa} = event.data;
    if (_res !== undefined) {
        await getResult(_res);
    }
};