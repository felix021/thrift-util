
enum EnumType {
    A = 1
    B = 2
}

struct Demo {
    1:  optional bool bool1,
    2:  optional double double1,
    3:  optional byte byte1,
    4:  optional i16 int16,
    5:  optional i32 int32,
    6:  optional i64 int64,
    7:  optional string string1,
    8:  optional Demo struct_demo,
    9:  optional map<string, i32> map_string_int32,
    10: optional map<string, Demo> map_string_demo,
    11: optional set<byte> set_byte,
    12: optional set<i32> set_i32,
    13: optional list<byte> list_byte,
    14: optional list<string> list_string,
    15: optional list<map<i64, Demo>> list_map_demo,
    16: optional EnumType enum1,
    // optional 16: uuid u,
}