
export default function() { 
    return {
      "view": function({attrs}) {
        return html`
        <div>
            <img src="/uploads/${attrs.section}.png" />
        </div>`
      }
    }
}