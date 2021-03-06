
export default function() { 
  return {
    "view": function({attrs}) {
      return html`
    <div class="bg-white">
      <div class="max-w-7xl mx-auto py-12 px-4 sm:px-6 md:py-16 lg:px-8 lg:py-20">
        <h2 class="text-3xl font-extrabold tracking-tight text-gray-900 sm:text-4xl">
          <span class="block">Ready to dive in?</span>
          <span class="block text-primary-600">Start your free trial today.</span>
        </h2>
        <div class="mt-8 flex">
          <div class="inline-flex rounded-md shadow">
            <a href="#" class="inline-flex items-center justify-center px-5 py-3 border border-transparent text-base font-medium rounded-md text-white bg-primary-600 hover:bg-primary-700">
              Get started
            </a>
          </div>
          <div class="ml-3 inline-flex">
            <a href="#" class="inline-flex items-center justify-center px-5 py-3 border border-transparent text-base font-medium rounded-md text-primary-700 bg-primary-100 hover:bg--200">
              Learn more
            </a>
          </div>
        </div>
      </div>
    </div>`
    }
  }
}