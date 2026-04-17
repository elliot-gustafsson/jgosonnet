// std.parseYaml('1\n---')

// std.objectRemoveKey({ foo: 1, bar: 2, baz: 3 }, 'foo')
// std.objectRemoveKey({ foo: 1, bar: 2, baz:: 3 }, 'foo').baz
// { x: 1 } + std.objectRemoveKey({ a: 1 } + { b: super.a }, 'a')

// { a: 1 } + std.objectRemoveKey({ b: super.a }, 'a')

// local object = {
//   foo: 'bar',
//   bar: self.foo,
//   baz: 1,
//   bazel: 1.42,
//   boom: -1,
//   bim: false,
//   bam: true,
//   blamo: {
//     cereal: [
//       '<>& fizbuzz',
//     ],

//     treats: [
//       {
//         name: 'chocolate',
//       },
//     ],
//   },
// };

// local array = [
//   'bar',
//   object.foo,
//   1,
//   1.42,
//   -1,
//   false,
//   true,
//   {
//     cereal: [
//       '<>& fizbuzz',
//     ],

//     treats: [
//       {
//         name: 'chocolate',
//       },
//     ],
//   },
// ];

// {
//   array: std.manifestJsonEx(array, '  '),
//   bool: std.manifestJsonEx(true, '   '),
//   'null': std.manifestJsonEx(null, '   '),
//   object: std.manifestJsonEx(object, '  '),
//   number: std.manifestJsonEx(42, '   '),
//   string: std.manifestJsonEx('foo', '   '),
// }


// {
//   a: std.manifestJson(1.42),
//   b: std.manifestJsonEx(1.42, indent=' ', newline='\n', key_val_sep=': '),
//   c: std.manifestJsonMinified(1.42),
//   d: std.toString({ a: 1.42 }),
//   e: std.toString(1.42),
//   f: std.manifestYamlDoc({ a: 1.42 }, indent_array_in_object=false, quote_keys=true),
//   g: std.manifestYamlDoc(1.42, indent_array_in_object=false, quote_keys=true),
// }

{
  a: std.manifestJson(9223372036854775807),
  b: std.manifestJsonEx(9223372036854775807, indent=' ', newline='\n', key_val_sep=': '),
  c: std.manifestJsonMinified(9223372036854775807),
  d: std.toString({ a: 9223372036854775807 }),
  e: std.toString(9223372036854775807),
  f: std.manifestYamlDoc({ a: 9223372036854775807 }, indent_array_in_object=false, quote_keys=true),
  g: std.manifestYamlDoc(9223372036854775807, indent_array_in_object=false, quote_keys=true),
}
