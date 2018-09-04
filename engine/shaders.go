package engine

var vertexShader = `
	#version 410

	// This shader will be called as many times as there are cells on the
	// screen and for single cell it will be called 6 times for every
	// vertex, because single cell consists of 2 triangles:
	//
	//   (0;1) *-* (1;1)
	//         |\|
	//   (0;0) *-* (1;0)
	//
	// gl_InstanceID will be incremented for every 6 vertices, since single
	// cell is one instance.

	// [input] in_Vertex: vertex coordinates for 6 vertices of single cell.
	in vec2 in_Vertex;

	// [input] in_Glyph: glyph coordinates in font, passed to fragment
	// shader (column, row).
	in ivec2 in_Glyph;

	// [input] in_Attrs: cell attributes, passed to fragment shader.
	in int in_Attrs;

	// [output] frag_Cell: (x, y) coordinates of single cell in pixels.
	flat out ivec2 frag_Cell;

	// [output] frag_Glyph: same as in_Glyph.
	flat out ivec2 frag_Glyph;

	// [output] frag_Attrs: same as in_Attrs.
	flat out int frag_Attrs;

	// [static] uni_ViewSize: viewport size in pixels.
	uniform ivec2 uni_ViewSize;

	// [static] uni_GlyphSize: glyph size in pixels.
	uniform ivec2 uni_GlyphSize;

	// We have local vertex coords for single cell and cell index.
	// From that we need to calculate vertex coords in world-space for
	// given cell.
	//
	// We do calculations in pixels for simplicity and then convert to
	// world-space coordinates.
	//
	// On-screen equvalent for world-space coordinates in vertex shader:
	// [-1; -1]: bottom left corner;
	// [ 0;  0]: center;
	// [ 1;  1]: top right corner;
	//
	void main() {
		int lineLength = uni_ViewSize.x / uni_GlyphSize.x;

		// Since gl_InstanceID is ranged from 0 to number of cells on
		// the screen, we need to map it to (row, column) coordinates.
		int row = gl_InstanceID / lineLength;
		int column = gl_InstanceID - row * lineLength; // mod() will not work

		frag_Glyph = in_Glyph;
		frag_Attrs = in_Attrs;
		frag_Cell = ivec2(column, row) * uni_GlyphSize;

		// After we have (row, column) coordinates, we can translate them
		// to world-space coordinates that are in range [-1; 1].
		//
		// frag_Cell — offset for drawing cell in pixels and
		// in_Vertex * uni_GlyphSize — coordinates of single vertex mapped
		// to pixel coordinates.
		vec2 coord = vec2(1, -1) * (
			2 * (frag_Cell + in_Vertex * uni_GlyphSize) / uni_ViewSize - 1
		);

		gl_Position = vec4(coord, 0, 1);
	}
`

var fragmentShader = `
	#version 410

	// [input] frag_Cell: pixel coordinates of fragment.
	flat in ivec2 frag_Cell;

	// [input] frag_Glyph: glyph coordinates in font.
	flat in ivec2 frag_Glyph;

	// [input] frag_Attrs: cell attributes in font.
	flat in int frag_Attrs;

	// [output] out_Color: color of current pixel.
	out vec4 out_Color;

	// [static] uni_ViewSize: viewport size in pixels.
	uniform ivec2 uni_ViewSize;

	// [static] uni_GlyphSize: glyph size in pixels.
	uniform ivec2 uni_GlyphSize;

	// [static] uni_Font: font texture.
	uniform sampler2D uni_Font;

	// We have pixel coordinates, font texture and top-left corner of
	// current cell.
	//
	// From that we need to get pixel color from font texture for given
	// glyph.
	//
	//                gl_FragCoord.y
	//                ^
	//                |
	//                |
	// (frag_Cell.xy) *--> gl_FragCoord.x
	//
	// Incoming gl_FragCoord are on-screen coordinates:
	// [0; 0]: bottom left corner (note: not *top* left);
	// [w; h]: top right corner;
	void main() {
		if (frag_Attrs == 0) {
			discard;
		}
		
		// We first calculate pixel coordinates in coordinate system of
		// current cell, e.g. from [0; 0] to [glyph width; glyph height].
		//
		// Because frag_Cell coorinates use top left corner as origin,
		// we need to inverse y-part.
		vec2 coord = (
			vec2(gl_FragCoord.x, uni_ViewSize.y - gl_FragCoord.y) -
			frag_Cell
		);

		// To obtain pixel color from font we do all calculations in pixels
		// and then convert them to texture-coords.
		//
		// frag_Glyph * uni_GlyphSize - origin of current glyph in texture,
		// coord - offset of pixel in glyph-local coordinates.
		out_Color = texture(
			uni_Font,
			(coord + frag_Glyph * uni_GlyphSize) / textureSize(uni_Font, 0)
		);
	}
`
