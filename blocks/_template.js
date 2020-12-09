import htm from 'https://unpkg.com/htm?module'
const html = htm.bind(h);

export default function() { 
    return html`
        <div>Hello world</div>`
}