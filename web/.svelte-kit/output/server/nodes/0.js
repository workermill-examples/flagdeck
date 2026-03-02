

export const index = 0;
let component_cache;
export const component = async () => component_cache ??= (await import('../entries/fallbacks/layout.svelte.js')).default;
export const imports = ["_app/immutable/nodes/0.Bg8vGak_.js","_app/immutable/chunks/D2YQVQxD.js","_app/immutable/chunks/C0Ec5M7Z.js","_app/immutable/chunks/Cht_CHgl.js"];
export const stylesheets = [];
export const fonts = [];
