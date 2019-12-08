require('fs').readFile('./input.txt', 'utf8', (_, c) =>
    console.log(((d, z) =>
      (x => x[1] * x[2])(Array((d.length / z)|0).fill(0)
        .map((_, i) => d.slice(i*z, i*z+z).reduce((a, v) => (a[v] = (a[v] || 0) + 1, a), {}))
        .reduce((a, v) => v[0] < a[0] ? v : a))
    )(c.split(''), 25*6)))

require('fs').readFile('./input.txt', 'utf8', (_, c) =>
    console.log(((d, w, h, z=w*h) =>
      (d => Array(h).fill(0).map((_, i) => d.slice(i*w, i*w+w).map(z => [' ', '*', ' '][z]).join('')).join('\n'))(
        Array((d.length / z)|0).fill(0)
          .map((_, i) => d.slice(i*z, i*z+z))
          .reduce((a, v) => a.map((w, i) => w === 2 ? v[i] : w))
      )
    )(c.split('').map(x=>+x), 25, 6)))
