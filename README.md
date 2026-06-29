# JavaScript parser

An extensible JavaScript parser mounted on top of [XJS](https://github.com/xjslang/xjs).

## And what makes it different from other parsers?

This one is "extensible". That is, we can easily add our own operators, expressions, and statements that aren't necessarily part of the JS standard.

For example, we could add the factorial operator. Syntactically, the expression `let x = 7!` should be incorrect, but after adding the operator, the parser would be able to recognize it and transform it into standard JS.

The possibilities are endless!

> [!NOTE]
> The project is currently in an early stage of development. I'll keep you updated :)
