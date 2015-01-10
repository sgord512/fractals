#version 330

#define MANDLEBROT 1
//#define CIRCLE 1
//#define NEWTON 1

const float maxDiffForFoundRoots = 1000.0;

uniform vec2 windowSize;
uniform vec2 viewOrigin;
uniform float viewZoom;
uniform int ul;
uniform float epsilon;
uniform int rootCount;
uniform vec3 defaultRootColor;
uniform vec2 rootBase[3];
uniform vec3 rootColor[3];
out vec4 outColor;

vec3 hsvToRgb(vec3 hsv) { 
  float c = hsv.z * hsv.y;
  float h = hsv.x * 6.0f;
  float x = c * (1 - abs(mod(h, 2) - 1));
  float m = hsv.z - c;
  vec3 rgb;
  if (h < 1) { 
    rgb = vec3(c, x, 0);
  } else if (h < 2) { 
    rgb = vec3(x, c, 0);
  } else if (h < 3) { 
    rgb = vec3(0, c, x);
  } else if (h < 4) { 
    rgb = vec3(0, x, c);
  } else if (h < 5) { 
    rgb = vec3(x, 0, c);
  } else {
    rgb = vec3(c, 0, x);    
  }
  return rgb + m;
}

vec2 conj(vec2 z) { 
  return vec2(z.x, -z.y);
}

float c_abs(vec2 z) { 
  return dot(z, z);
}

vec2 reciprocal(vec2 z) { 
  return conj(z) / c_abs(z);
}

vec2 evaluateAt(vec2 z) { 
  vec2 y = vec2(1.0, 0.0);
  for (int i = 0; i < rootCount; i++) {
    y *= z - rootBase[i];
  }
  return y;
}

#ifdef NEWTON
  vec2 evaluateDeltaAt(vec2 z) { 
    vec2 total = vec2(0.0);
    for (int i = 0; i < rootCount; i++) { 
      total += reciprocal(z - rootBase[i]);
    }
    return reciprocal(total);
  }

  int closestRoot(vec2 z) {
    int ix = -1;
    float minDiff = 1.0 / 0.0;
    for (int i = 0; i < rootCount; i++) { 
      float currDiff = c_abs(z - rootBase[i]);
      if (currDiff < minDiff) { 
	ix = i;
	minDiff = currDiff;
      }
    }
    if (minDiff > maxDiffForFoundRoots) {
      return -1;
    }
    return ix;
  }

  vec2 newtonRaphson(vec2 z) { 
    vec2 y = evaluateAt(z);
    for(int count = 0; count < ul; count++) { 
      z = z - evaluateDeltaAt(z);
      y = evaluateAt(z);
      if (c_abs(y) < epsilon / viewZoom) { 
	return vec2(count, closestRoot(z));
      }
    }
    return vec2(ul, closestRoot(z));
  }
#endif // NEWTON

#ifdef CIRCLE
  vec2 withinCircle(vec2 z) { 
    if (c_abs(z) <= 2.0) { 
      return vec2(0.0);
    } else {
      return vec2(ul, -1.0);
    }
  }
#endif // CIRCLE

#ifdef MANDLEBROT
  vec2 mandlebrot(vec2 c) {
    vec2 z = c;
    for (int count = 0; count < ul; count++) { 
      z = vec2(pow(z.x, 2) - pow(z.y, 2), 2 * z.x * z.y) + c;
      if (c_abs(z) >= 2.0) { 
	return vec2(count, 0.0);
      }
    }
    return vec2(ul, -1.0);
  }
#endif // MANDLEBROT

void main()
{
    vec2 windowCoord = gl_FragCoord.xy / windowSize;
    vec2 windowAspect = vec2(1.0);
    if (windowSize.x >= windowSize.y) {       
      windowAspect.x = windowSize.x / windowSize.y;
    } else {
      windowAspect.y = windowSize.y / windowSize.x;
    }
    vec2 windowOriginOffset = windowCoord - vec2(0.5);
    vec2 viewOriginOffset = windowOriginOffset * windowAspect / vec2(viewZoom);
    vec2 location = viewOriginOffset + viewOrigin;    
    
    #ifdef NEWTON
      vec2 results = newtonRaphson(location);
    #elif MANDLEBROT
      vec2 results = mandlebrot(location);
    #elif CIRCLE
      vec2 results = withinCircle(location)
    #endif
    
    float count = results.x;
    int rootIx = int(results.y);
    vec3 hue = defaultRootColor;
    if (rootIx != -1) { 
      hue = rootColor[rootIx];
    }

    //hue = hsvToRgb(vec3(count / ul, 1.0, 1.0));

    outColor = vec4(hue, float(results.x / ul));
}
