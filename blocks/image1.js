import htm from 'https://unpkg.com/htm?module'
const html = htm.bind(h);

export default function({attrs}) { 
    var section = attrs.section;
    return html`
        <div>
            <img src="/uploads/${section}.png" />
        </div>`
}