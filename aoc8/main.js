require('fs').readFile('./input.txt', 'utf8', (_, c) =>
    console.log(((d, z) =>
      [Array(d.length / z).fill(0)
        .map((_, i) => d.slice(i*z, i*z+z).reduce((a, v) => (a[v] = (a[v] || 0) + 1, a), {}))
        .reduce((a, v) => v[0] < a[0] ? v : a)].map(x => x[1] * x[2])[0]
    )(c.trim().split(''), 25*6)))
