let back_url = 'https://e-sketch.ru'

const isLocalhost = Boolean(
    window.location.hostname === 'localhost' ||
    window.location.hostname === '[::1]' ||
    window.location.hostname.match(
        /^127(?:\.(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}$/
    )
);

if (isLocalhost) {
    back_url = window.location.protocol + '//' + window.location.hostname + ':' + window.location.port
}

export const Config = {
    back_url: back_url,
}