{
	description = "Label utilities";

	inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-25.11";

		flake-utils.url = "github:numtide/flake-utils";

		templ.url = "github:a-h/templ/v0.3.960";
		templ.inputs.nixpkgs.follows = "nixpkgs";
	};

	outputs = { self, nixpkgs, flake-utils, ... }@inputs: flake-utils.lib.eachDefaultSystem (system:
		let
			pkgs = import nixpkgs { inherit system; };
			templ = inputs.templ.packages.${system}.templ;
		in
		{
			packages.textlabel = pkgs.buildGoModule {
				pname = "textlabel";
				version = "0.1.0";

				src = ./.;

				vendorHash = null;
				subPackages = [ "./textlabel" ];
			};

			packages.printserver = pkgs.buildGoModule {
				pname = "printserver";
				version = "0.1.0";

				src = ./.;

				vendorHash = null;
				subPackages = [ "./printserver" ];
			};

			packages.webprint = pkgs.buildGoModule rec {
				pname = "webprint";
				version = "0.1.0";

				src = ./.;

				preBuild = ''
					${templ}/bin/templ generate
					${pkgs.coreutils}/bin/mkdir -p $out/bin
					${pkgs.coreutils}/bin/mkdir -p $out/share/${pname}/static/wasm
				'';

				serveAddr = "0.0.0.0:8080";

				buildPhase = ''
					${pkgs.go}/bin/go build -o $out/bin/webprint -ldflags="-s -w -X 'main.StaticDir=$out/share/${pname}/static' -X 'main.ServeAddr=${serveAddr}'" ./webprint/server
					GOOS=js GOARCH=wasm ${pkgs.go}/bin/go build -o $out/share/${pname}/static/wasm/main.wasm ./webprint/wasm
				'';

				installPhase = ''
					${pkgs.coreutils}/bin/cp -a ${./webprint/static}/. $out/share/${pname}/static
					${pkgs.coreutils}/bin/cp "$(${pkgs.go}/bin/go env GOROOT)/lib/wasm/wasm_exec.js" $out/share/${pname}/static/wasm/wasm_exec.js
				'';

				vendorHash = null;
			};

			devShells.default = pkgs.mkShell {
				buildInputs = with pkgs; [
					go
					templ
				];
			};

			devShells.webprint = pkgs.mkShell {
				buildInputs = with pkgs; [
					go
					templ
				];
				shellHook = ''
					${pkgs.coreutils}/bin/mkdir -p ./webprint/static/wasm
					${pkgs.coreutils}/bin/cp "$(${pkgs.go}/bin/go env GOROOT)/lib/wasm/wasm_exec.js" ./webprint/static/wasm/wasm_exec.js
					runwebprint() {
						GOOS=js GOARCH=wasm ${pkgs.go}/bin/go build -o ./webprint/static/wasm/main.wasm ./webprint/wasm
						${pkgs.go}/bin/go run -ldflags="-X main.ServeAddr=:8080 -X main.StaticDir=./webprint/static" ./webprint/server ws://127.0.0.1:6245
					}
					export -f runwebprint
					${pkgs.coreutils}/bin/echo "Watching webprint..."
					${templ}/bin/templ generate --watch --cmd="${pkgs.bash}/bin/bash -c runwebprint" --proxy=http://127.0.0.1:8080
				'';
			};
		}
	);
}
