# aws-iam-conditionals-as-cel

This repository contains a few Golang example tests that showcase [cel-go](https://github.com/google/cel-go) by showing how the following AWS IAM policy JSON could be expressed with CEL:

```json
[
    {"Condition" : { "StringEqualsIgnoreCase" : { "aws:username" : "johndoe" }}},
    {"Condition": { "NumericLessThanEquals": {"aws:MultiFactorAuthAge": "3600"}}},
    {"Condition": { "ForAnyValue:StringEquals": {
        "aws:TagKeys": ["foo", "bar"]
    }}}
]
```

Final result:

```go
user.Name.lowerAscii() == "johndoe" &&
    user.multiFactorAuthAge() <= uint(3600) &&
    user.Tags.exists(tag, tag.Key == "foo" || tag.Key == "bar")
```
