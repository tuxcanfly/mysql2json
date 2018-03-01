mysql2json:
===========


`mysql2json` converts a mysql dump to JSON

Usage:
------

`mysql2json <input/file/path>`

where input file is a mysql dump of the form:

    INSERT INTO `departments` VALUES
    ('d001','Marketing'),
    ('d002','Finance');

the JSON will be written to stdout as:

    {"TableName":"departments","Values":{"d001":["Marketing"],"d002":["Finance"]}}
