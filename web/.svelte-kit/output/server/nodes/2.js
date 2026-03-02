

export const index = 2;
let component_cache;
export const component = async () => component_cache ??= (await import('../entries/pages/_page.svelte.js')).default;
export const imports = ["_app/immutable/nodes/2.DdYz2m8Y.js","_app/immutable/chunks/D2YQVQxD.js","_app/immutable/chunks/C0Ec5M7Z.js","_app/immutable/chunks/DcEZgWxx.js"];
export const stylesheets = ["_app/immutable/assets/2.DoG9X-Sa.css"];
export const fonts = [];
