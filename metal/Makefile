METAL := xcrun -sdk macosx metal
METALLIB := xcrun -sdk macosx metallib
SHADER_SRC := ../internal/gpu/eml_metal_shaders.metal

.PHONY: all clean

all: eml_metal_shaders.metallib

eml_metal_shaders.air: $(SHADER_SRC)
	$(METAL) -fno-fast-math -Wall -Wextra -c $(SHADER_SRC) -o $@

eml_metal_shaders.metallib: eml_metal_shaders.air
	$(METALLIB) $< -o $@

clean:
	rm -f eml_metal_shaders.air eml_metal_shaders.metallib
