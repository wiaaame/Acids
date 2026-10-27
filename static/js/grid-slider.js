(function () {
    const minEl = document.getElementById('rangeMin');
    const maxEl = document.getElementById('rangeMax');
    const fill  = document.getElementById('rangeFill');
    const bMin  = document.getElementById('bubbleMin');
    const bMax  = document.getElementById('bubbleMax');
    if (!minEl || !maxEl || !fill) return;

    const MIN = +minEl.min, MAX = +maxEl.max;
    const pct = v => ((v - MIN) / (MAX - MIN)) * 100;

    function render() {
        let a = +minEl.value, b = +maxEl.value;
        if (a > b) { [a, b] = [b, a]; }
        fill.style.left  = pct(a) + '%';
        fill.style.width = (pct(b) - pct(a)) + '%';
        bMin.textContent = a;
        bMax.textContent = b;
        bMin.style.left  = pct(a) + '%';
        bMax.style.left  = pct(b) + '%';
    }

    minEl.addEventListener('input', render);
    maxEl.addEventListener('input', render);
    render();
})();