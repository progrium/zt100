
export default function() { 
    return {
      "view": function({attrs}) {
        return html`
      <div class="relative bg-gray-50 pt-16 pb-20 px-4 sm:px-6 lg:pt-12 lg:pb-14 lg:px-8">
        <div class="relative bg-white max-w-5xl mx-auto p-12">
            <h2 class="text-xl">Consent</h2>
            <textarea class="w-full h-24">
                Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.
            </textarea>
            <form method="POST" class="w-full text-right"><input class="ml-8 whitespace-nowrap inline-flex items-center justify-center px-4 py-2 border border-transparent rounded-md shadow-sm text-base font-medium text-white bg-primary-500 hover:bg-primary-500" type="submit" value="I Agree" /></form>
        </div>
      </div>`
      }
    }
}