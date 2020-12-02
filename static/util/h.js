import { Style } from "./style.js";

var _h = {};
export function h(tag, attrs, ...children) {
    //console.log(tag, typeof tag === "function", !isClass(tag), _h[tag] == undefined)
    if (typeof tag === "function" && !isClass(tag) && _h[tag] == undefined) {
        //console.log(tag.prototype.constructor.name);
        _h[tag] = wrap(tag)
    }
    if (_h[tag]) {
        tag = _h[tag];
    }
    return m(tag, attrs, children);
}
h.redraw = m.redraw;
h.mount = m.mount;
h.render = m.render;
h.trust = m.trust;
window.h = h;

function wrap(v) {
    return {
        view: function(input) {
            let style = Style.from(v);
            
            input.hooks = {};
            input.style = style;
            input.vnode = input;
            input.attrs = attrProxy(input.attrs);

            let output = v(input);
            
            applyAttrs(style, input.attrs);
            applyHooks(output, input.hooks);
            applyStyle(output, style);
            applyEvents(output, input.attrs);
            applyId(output, input.attrs);
            applyData(output, input.attrs);
            
            return output;
        }
    }  
}

function applyData(vnode, attrs) {
    for (let attr in attrs) {
        if (attr.startsWith("data-")) {
            vnode.attrs[attr] = attrs[attr];
        }
    }
}

function applyId(vnode, attrs) {
    if (attrs.id) {
        vnode.attrs.id = attrs.id;
    }
}

function applyAttrs(style, attrs) {
    if (attrs.class) {
        style.add(attrs.class);
    }
    if (attrs.style) {
        style.add(attrs.style);
    }
}

function applyStyle(vnode, style) {
    if (vnode.attrs.class === undefined) {
        vnode.attrs.class = style.class();
    }
    if (vnode.attrs.style === undefined) {
        vnode.attrs.style = style.style();
    }
}

function applyHooks(vnode, hooks) {
    vnode.attrs = Object.assign(vnode.attrs||{}, hooks);
}

function applyEvents(vnode, attrs) {
    for (let attr in attrs) {
        if (attr.startsWith("on") && !attrs._used.has(attr)) {
            vnode.attrs[attr] = attrs[attr];
        }
    }
}

function isClass(obj) {
    if (!obj.prototype || !obj.prototype.constructor) return false;
    return obj.prototype.constructor.toString().startsWith("class ");
}

function attrProxy(attrs) {
    return new Proxy(attrs, {
        get: function (target, prop, receiver) {
            if (!this.used) {
                this.used = new Set();
            }
            if (prop === "_used") {
                return this.used;
            }
            if (prop.startsWith("on")) {
                this.used.add(prop);
            }
            return Reflect.get(...arguments);
        },
    })
}
