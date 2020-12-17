
export default function() { 
  return {
    "view": function({attrs}) {
      return html`
    <main class="mt-16 sm:mt-24">
      <div class="mx-auto max-w-7xl">
        <div class="lg:grid lg:grid-cols-12 lg:gap-8">
          <div class="px-4 sm:px-6 sm:text-center md:max-w-2xl md:mx-auto lg:col-span-6 lg:text-left lg:flex lg:items-center">
            <div>
              <a href="#" class="inline-flex items-center text-white bg-gray-900 rounded-full p-1 pr-2 sm:text-base lg:text-sm xl:text-base hover:text-gray-200">
                <span class="px-3 py-0.5 text-white text-xs font-semibold leading-5 uppercase tracking-wide bg-primary-500 rounded-full">We're hiring</span>
                <span class="ml-4 text-sm">Visit our careers page</span>
                
                <svg class="ml-2 w-5 h-5 text-gray-500" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
                  <path fill-rule="evenodd" d="M7.293 14.707a1 1 0 010-1.414L10.586 10 7.293 6.707a1 1 0 011.414-1.414l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0z" clip-rule="evenodd" />
                </svg>
              </a>
              <h1 class="mt-4 text-4xl tracking-tight font-extrabold text-white sm:mt-5 sm:leading-none lg:mt-6 lg:text-5xl xl:text-6xl">
                <span class="md:block">Data to enrich your</span>
                <span class="text-primary-400 md:block">online business</span>
              </h1>
              <p class="mt-3 text-base text-gray-300 sm:mt-5 sm:text-xl lg:text-lg xl:text-xl">
                Anim aute id magna aliqua ad ad non deserunt sunt. Qui irure qui lorem cupidatat commodo. Elit sunt amet fugiat veniam occaecat fugiat aliqua ad ad non deserunt sunt.
              </p>
              <p class="mt-8 text-sm text-white uppercase tracking-wide font-semibold sm:mt-10">Used by</p>
              <div class="mt-5 w-full sm:mx-auto sm:max-w-lg lg:ml-0">
                <div class="flex flex-wrap items-start justify-between">
                  <div class="flex justify-center px-1">
                    <img class="h-9 sm:h-10" src="https://tailwindui.com/img/logos/tuple-logo-gray-400.svg" alt="Tuple" />
                  </div>
                  <div class="flex justify-center px-1">
                    <img class="h-9 sm:h-10" src="https://tailwindui.com/img/logos/workcation-logo-gray-400.svg" alt="Workcation" />
                  </div>
                  <div class="flex justify-center px-1">
                    <img class="h-9 sm:h-10" src="https://tailwindui.com/img/logos/statickit-logo-gray-400.svg" alt="StaticKit" />
                  </div>
                </div>
              </div>
            </div>
          </div>
          <div class="mt-16 sm:mt-24 lg:mt-0 lg:col-span-6">
            <div class="bg-white sm:max-w-md sm:w-full sm:mx-auto sm:rounded-lg sm:overflow-hidden">
              <div class="px-4 py-8 sm:px-10">
                <div>
                  <p class="text-sm font-medium text-gray-700">
                    Sign in with
                  </p>

                  <div class="mt-1 grid grid-cols-3 gap-3">
                    <div>
                      <a href="#" class="w-full inline-flex justify-center py-2 px-4 border border-gray-300 rounded-md shadow-sm bg-white text-sm font-medium text-gray-500 hover:bg-gray-50">
                        <span class="sr-only">Sign in with Facebook</span>
                        <svg class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20" aria-hidden="true">
                          <path fill-rule="evenodd" d="M20 10c0-5.523-4.477-10-10-10S0 4.477 0 10c0 4.991 3.657 9.128 8.438 9.878v-6.987h-2.54V10h2.54V7.797c0-2.506 1.492-3.89 3.777-3.89 1.094 0 2.238.195 2.238.195v2.46h-1.26c-1.243 0-1.63.771-1.63 1.562V10h2.773l-.443 2.89h-2.33v6.988C16.343 19.128 20 14.991 20 10z" clip-rule="evenodd" />
                        </svg>
                      </a>
                    </div>

                    <div>
                      <a href="#" class="w-full inline-flex justify-center py-2 px-4 border border-gray-300 rounded-md shadow-sm bg-white text-sm font-medium text-gray-500 hover:bg-gray-50">
                        <span class="sr-only">Sign in with Twitter</span>
                        <svg class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20" aria-hidden="true">
                          <path d="M6.29 18.251c7.547 0 11.675-6.253 11.675-11.675 0-.178 0-.355-.012-.53A8.348 8.348 0 0020 3.92a8.19 8.19 0 01-2.357.646 4.118 4.118 0 001.804-2.27 8.224 8.224 0 01-2.605.996 4.107 4.107 0 00-6.993 3.743 11.65 11.65 0 01-8.457-4.287 4.106 4.106 0 001.27 5.477A4.073 4.073 0 01.8 7.713v.052a4.105 4.105 0 003.292 4.022 4.095 4.095 0 01-1.853.07 4.108 4.108 0 003.834 2.85A8.233 8.233 0 010 16.407a11.616 11.616 0 006.29 1.84" />
                        </svg>
                      </a>
                    </div>

                    <div>
                      <a href="#" class="w-full inline-flex justify-center py-2 px-4 border border-gray-300 rounded-md shadow-sm bg-white text-sm font-medium text-gray-500 hover:bg-gray-50">
                        <span class="sr-only">Sign in with GitHub</span>
                        <svg class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20" aria-hidden="true">
                          <path fill-rule="evenodd" d="M10 0C4.477 0 0 4.484 0 10.017c0 4.425 2.865 8.18 6.839 9.504.5.092.682-.217.682-.483 0-.237-.008-.868-.013-1.703-2.782.605-3.369-1.343-3.369-1.343-.454-1.158-1.11-1.466-1.11-1.466-.908-.62.069-.608.069-.608 1.003.07 1.531 1.032 1.531 1.032.892 1.53 2.341 1.088 2.91.832.092-.647.35-1.088.636-1.338-2.22-.253-4.555-1.113-4.555-4.951 0-1.093.39-1.988 1.029-2.688-.103-.253-.446-1.272.098-2.65 0 0 .84-.27 2.75 1.026A9.564 9.564 0 0110 4.844c.85.004 1.705.115 2.504.337 1.909-1.296 2.747-1.027 2.747-1.027.546 1.379.203 2.398.1 2.651.64.7 1.028 1.595 1.028 2.688 0 3.848-2.339 4.695-4.566 4.942.359.31.678.921.678 1.856 0 1.338-.012 2.419-.012 2.747 0 .268.18.58.688.482A10.019 10.019 0 0020 10.017C20 4.484 15.522 0 10 0z" clip-rule="evenodd" />
                        </svg>
                      </a>
                    </div>
                  </div>
                </div>

                <div class="mt-6 relative">
                  <div class="absolute inset-0 flex items-center" aria-hidden="true">
                    <div class="w-full border-t border-gray-300"></div>
                  </div>
                  <div class="relative flex justify-center text-sm">
                    <span class="px-2 bg-white text-gray-500">
                      Or
                    </span>
                  </div>
                </div>

                <div class="mt-6">
                  <form action="#" method="POST" class="space-y-6">
                    <div>
                      <label for="name" class="sr-only">Full name</label>
                      <input type="text" id="name" placeholder="Full name" required class="block w-full shadow-sm focus:ring-primary-500 focus:border-primary-500 sm:text-sm border-gray-300 rounded-md" />
                    </div>

                    <div>
                      <label for="mobile-or-email" class="sr-only">Mobile number or email</label>
                      <input type="text" id="mobile-or-email" placeholder="Mobile number or email" required class="block w-full shadow-sm focus:ring-primary-500 focus:border-primary-500 sm:text-sm border-gray-300 rounded-md" />
                    </div>

                    <div>
                      <label for="password" class="sr-only">Password</label>
                      <input id="password" type="password" placeholder="Password" required class="block w-full shadow-sm focus:ring-primary-500 focus:border-primary-500 sm:text-sm border-gray-300 rounded-md" />
                    </div>

                    <div>
                      <button type="submit" class="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-primary-600 hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500">
                        Create your account
                      </button>
                    </div>
                  </form>
                </div>
              </div>
              <div class="px-4 py-6 bg-gray-50 border-t-2 border-gray-200 sm:px-10">
                <p class="text-xs leading-5 text-gray-500">By signing up, you agree to our <a href="#" class="font-medium text-gray-900 hover:underline">Terms</a>, <a href="#" class="font-medium text-gray-900 hover:underline">Data Policy</a> and <a href="#" class="font-medium text-gray-900 hover:underline">Cookies Policy</a>.</p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </main>`
    }
  }
}