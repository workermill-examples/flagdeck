

export const index = 2;
let component_cache;
export const component = async () => component_cache ??= (await import('../entries/pages/_page.svelte.js')).default;
export const imports = ["_app/immutable/nodes/2.BWcgz4YC.js","_app/immutable/chunks/wYGahFLe.js","_app/immutable/chunks/BDMaAkqK.js","_app/immutable/chunks/FMd03m4b.js"];
export const stylesheets = ["_app/immutable/assets/2.ClJf9tOG.css"];
export const fonts = [];
