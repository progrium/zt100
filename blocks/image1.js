
export default function() { 
    return {
      "view": function({attrs}) {
        return html`
        <div>
            <a href="login"></a><img src="/uploads/${attrs.section}.png" /></a>
        </div>`
      }
    }
}