export const manifest = (() => {
function __memo(fn) {
	let value;
	return () => value ??= (value = fn());
}

return {
	appDir: "_app",
	appPath: "_app",
	assets: new Set([]),
	mimeTypes: {},
	_: {
		client: {start:"_app/immutable/entry/start.D03I7dDP.js",app:"_app/immutable/entry/app.DDqGY-jF.js",imports:["_app/immutable/entry/start.D03I7dDP.js","_app/immutable/chunks/DBvXese6.js","_app/immutable/chunks/BDMaAkqK.js","_app/immutable/entry/app.DDqGY-jF.js","_app/immutable/chunks/BDMaAkqK.js","_app/immutable/chunks/FMd03m4b.js","_app/immutable/chunks/wYGahFLe.js","_app/immutable/chunks/DD2XqYrE.js"],stylesheets:[],fonts:[],uses_env_dynamic_public:false},
		nodes: [
			__memo(() => import('./nodes/0.js')),
			__memo(() => import('./nodes/1.js')),
			__memo(() => import('./nodes/2.js'))
		],
		remotes: {
			
		},
		routes: [
			{
				id: "/",
				pattern: /^\/$/,
				params: [],
				page: { layouts: [0,], errors: [1,], leaf: 2 },
				endpoint: null
			}
		],
		prerendered_routes: new Set([]),
		matchers: async () => {
			
			return {  };
		},
		server_assets: {}
	}
}
})();
