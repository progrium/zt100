export function Editable({attrs,style,children,hooks}) {
    var overrides = attrs.overrides || {};

    hooks.oncreate = (v) => {
        new MediumEditor(v.dom);
        v.dom.addEventListener("input", () => {
            debouncer(500, (data) => {
                if (data.substr(0, 3) === "<p>") {
                    data = data.substr(3, data.length-7);
                }
                m.request({
                    method: "POST",
                    url: "/c/zt100.override-text",
                    body: {
                        Text: data, 
                        Tenant: window.config.Tenant, 
                        Page: window.config.Page, 
                        App: window.config.App, 
                        Section: attrs.section,
                        Key: attrs.key 
                    },
                })
            }, v.dom.innerHTML);
        })
    }
    
    style.add("border-dashed border-2 border-gray-500 border-opacity-0 hover:border-opacity-100 focus:border-opacity-0");
    if (overrides[attrs.key] !== undefined) {
        return h("div", {}, h.trust(overrides[attrs.key]));
    }
    return h("div", {}, children);
}

let timeout = undefined;
const debouncer = (wait, fn, ...args) => {
    clearTimeout(timeout);
    const later = () => {clearTimeout(timeout); fn(...args);}
    timeout = (wait === 0) ? later() : setTimeout(later, wait);
};