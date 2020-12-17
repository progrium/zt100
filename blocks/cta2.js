
export default function() { 
  return {
    "view": function({attrs}) {
      return html`
    <div class="bg-primary-700">
      <div class="max-w-2xl mx-auto text-center py-16 px-4 sm:py-20 sm:px-6 lg:px-8">
        <h2 class="text-3xl font-extrabold text-white sm:text-4xl">
          <span class="block">Boost your productivity.</span>
          <span class="block">Start using Workflow today.</span>
        </h2>
        <p class="mt-4 text-lg leading-6 text-primary-200">Ac euismod vel sit maecenas id pellentesque eu sed consectetur. Malesuada adipiscing sagittis vel nulla nec.</p>
        <a href="#" class="mt-8 w-full inline-flex items-center justify-center px-5 py-3 border border-transparent text-base font-medium rounded-md text-primary-600 bg-white hover:bg-primary-50 sm:w-auto">
          Sign up for free
        </a>
      </div>
    </div>`
    }
  }
}