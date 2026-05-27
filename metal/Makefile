METAL := xcrun -sdk macosx metal
METALLIB := xcrun -sdk macosx metallib

SHADERS_SRC := eml_metal_shaders.metal
SHADERS_AIR := eml_metal_shaders.air
SHADERS_LIB := eml_metal_shaders.metallib

.PHONY: all clean

all: $(SHADERS_LIB)

$(SHADERS_AIR): $(SHADERS_SRC)
	$(METAL) -std=metal3.0 -o $@ -c $<

$(SHADERS_LIB): $(SHADERS_AIR)
	$(METALLIB) -o $@ $<

clean:
	rm -f $(SHADERS_AIR) $(SHADERS_LIB)
