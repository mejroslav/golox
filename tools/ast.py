# Create code for AST classes in Go.

import argparse

def generate_ast(output_dir: str) -> str:

    define_ast(output_dir, "Expr", [
        "Binary   : Left Expr, Operator *Token, Right Expr",
        "Grouping : Expression Expr",
        "Literal  : Value any",
        "Unary    : Operator *Token, Right Expr",
    ])

def define_ast(output_dir: str, file: str, types: list[str]) -> None:
    path = f"{output_dir}/{file}.go"
    with open(path, "w") as f:
        f.write("package golox\n\n")
        f.write("type " + file + " interface {\n")
        f.write("\tAccept(visitor " + file + "Visitor) any\n")
        f.write("}\n\n")

        # Visitor interface
        f.write("type " + file + "Visitor interface {\n")
        for type_def in types:
            class_name = type_def.split(":")[0].strip()
            f.write(f"\tVisit{class_name}{file}({file.lower()} *{class_name}) any\n")
        f.write("}\n\n")

        # AST classes
        for type_def in types:
            class_name, fields = type_def.split(":")
            class_name = class_name.strip()
            fields = fields.strip()

            f.write(f"type {class_name} struct {{\n")
            for field in fields.split(","):
                name, type_ = field.strip().split(" ")
                f.write(f"\t{name.capitalize()} {type_}\n")
            f.write("}\n\n")

            # Accept method
            f.write(f"func (node *{class_name}) Accept(visitor {file}Visitor) any {{\n")
            f.write(f"\treturn visitor.Visit{class_name}{file}(node)\n")
            f.write("}\n\n")


if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument(
        "--output",
        "-o",
        type=str,
        required=True,
        help="The output directory to write the generated AST code to.",
    )
    args = parser.parse_args()

    ast = generate_ast(args.output)
