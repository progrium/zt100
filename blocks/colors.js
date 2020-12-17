
export default function() { 
  return {
    "view": function() {
      return html`
        <div class="flex flex-row w-full text-center">
          <div class="bg-primary-200 flex-grow">200</div>
          <div class="bg-primary-300 flex-grow">300</div>
          <div class="bg-primary-400 flex-grow">400</div>
          <div class="bg-primary-500 flex-grow">500</div>
          <div class="bg-primary-600 flex-grow">600</div>
          <div class="bg-primary-700 flex-grow">700</div>
          <div class="bg-primary-800 flex-grow">800</div>
        </div>`
    }
  }
}
