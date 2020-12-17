
export default function() { 
    return {
      "view": function({attrs}) {
        return html`
      <div class="relative bg-gray-50 pt-16 pb-20 px-4 sm:px-6 lg:pt-12 lg:pb-14 lg:px-8">
        <div class="relative bg-white max-w-5xl mx-auto p-12">
          <pre>
              ${data.Contrib.Auth.IsAuthenticated && JSON.stringify(data.Contrib.Auth.Profile, null, 2)}
          </pre>
        </div>
      </div>`
      }
    }
}