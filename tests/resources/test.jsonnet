// std.parseYaml('1\n---')

// std.objectRemoveKey({ foo: 1, bar: 2, baz: 3 }, 'foo')
// std.objectRemoveKey({ foo: 1, bar: 2, baz:: 3 }, 'foo').baz
// { x: 1 } + std.objectRemoveKey({ a: 1 } + { b: super.a }, 'a')

// { a: 1 } + std.objectRemoveKey({ b: super.a }, 'a')

// std.objectRemoveKey({ a: 1 } + { b: super.a }, 'a')

local bare_yaml_quoted = {
  '685230': 'canonical',
  '+685_230': 'decimal',
  '02472256': 'octal',
  '-1_0': 'negative integer',
  '-0.1_0_0': 'negative float',
  '0x_0A_74_AE': 'hexadecimal',
  '-0x_0A_74_AE': 'negative hexadecimal',
  '0b1010_0111_0100_1010_1110': 'binary',
  '-0b1010_0111_0100_1010_1110': 'binary',
  '190:20:30': 'sexagesimal',
  '-190:20:30': 'negative sexagesimal',
  '6.8523015e+5': 'canonical',
  '6.8523015e-5': 'canonical',
  '-6.8523015e+5': 'negative canonical',
  '685.230_15e+03': 'exponential',
  '-685.230_15e+03': 'negative exponential',
  '-685.230_15e-03': 'negative w/ negative exponential',
  '-685.230_15E-03': 'negative w/ negative exponential',
  '685_230.15': 'fixed',
  '-685_230.15': 'negative fixed',
  '190:20:30.15': 'sexagesimal',
  '-190:20:30.15': 'negative sexagesimal',
  '-.inf': 'negative infinity',
  '.inf': 'positive infinity',
  '+.inf': 'positive infinity',
  '.NaN': 'not a number',
  y: 'boolean true',
  yes: 'boolean true',
  Yes: 'boolean true',
  True: 'boolean true',
  'true': 'boolean true',
  on: 'boolean true',
  On: 'boolean true',
  NO: 'boolean false',
  n: 'boolean false',
  N: 'boolean false',
  off: 'boolean false',
  OFF: 'boolean false',
  'null': 'null word',
  NULL: 'null word capital',
  Null: 'null word',
  '~': 'null key',
  '': 'empty key',
  '-': 'invalid bare key',
  '---': 'triple dash key',
  '2001-12-15T02:59:43.1Z': 'canonical',
  '2001-12-14t21:59:43.10-05:00': 'valid iso8601',
  '2001-12-14 21:59:43.10 -5': 'space separated',
  '2001-12-15 2:59:43.10': 'no time zone (Z)',
  '2002-12-14': 'date',
};
local bare_yaml_unquoted = {
  '0X_0a_74_ae': 'BARE_KEY',
  '__-0X_0a_74_ae': 'BARE_KEY',
  '-0B1010_0111_0100_1010_1110': 'BARE_KEY',
  '__-0B1010_0111_0100_1010_1110': 'BARE_KEY',
  x: 'BARE_KEY',
  b: 'BARE_KEY',
  just_letters_underscores: 'BARE_KEY',
  'just-letters-dashes': 'BARE_KEY',
  'jsonnet.org/k8s-label-like': 'BARE_KEY',
  '192.168.0.1': 'BARE_KEY',
  '1-234-567-8901': 'BARE_KEY',
};
local bare_yaml_test = bare_yaml_quoted + bare_yaml_unquoted;

local x = {
  'jsonnet.org/k8s-label-like': 'asdf',
  '2001-12-15T02:59:43.1Z': 'canonical',
  '2001-12-14t21:59:43.10-05:00': 'valid iso8601',
  '2001-12-14 21:59:43.10 -5': 'space separated',
  '2001-12-15 2:59:43.10': 'no time zone (Z)',
  '2002-12-14': 'date',
};

// x

std.manifestYamlDoc(x, quote_keys=false)

// {
//   '2002-12-14': 'asdf',
// }
