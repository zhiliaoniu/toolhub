template <typename T, size_t N> char (&ArraySizeHelper(T (&array)[N]))[N];
#define ABSL_ARRAYSIZE(array) (sizeof::ArraySizeHelper(array))
