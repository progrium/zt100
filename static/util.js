export function HasFeature(flag) {
    return window.data.Demo.Features.includes(flag);
}

export function Editable({attrs,style,children,hooks}) {
    var text = attrs.text || {};
    var key = attrs.key; // key
    var id = attrs.id; // block id
    

    hooks.oncreate = (v) => {
        new MediumEditor(v.dom);
        v.dom.addEventListener("input", () => {
            debouncer(500, (data) => {
                if (data.substr(0, 3) === "<p>") {
                    data = data.substr(3, data.length-7);
                }
                m.request({
                    method: "POST",
                    url: "/cmd/zt100.block.set-text",
                    body: {
                        Text: data, 
                        Key: key,
                        BlockID: id, 
                        PageID: window.config.Page.OID,
                        
                    },
                })
            }, v.dom.innerHTML);
        })
    }
    
    style.add("border-dashed border-2 border-gray-500 border-opacity-0 hover:border-opacity-100 focus:border-opacity-0");
    if (text[key] !== undefined) {
        return h("div", {}, h.trust(text[key]));
    }
    return h("div", {}, children);
}

let timeout = undefined;
const debouncer = (wait, fn, ...args) => {
    clearTimeout(timeout);
    const later = () => {clearTimeout(timeout); fn(...args);}
    timeout = (wait === 0) ? later() : setTimeout(later, wait);
};