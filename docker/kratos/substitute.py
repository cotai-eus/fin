#!/usr/bin/env python3
import sys
import os

template_file = sys.argv[1]
output_file = sys.argv[2]

vars_to_sub = {
    'KRATOS_DSN': os.environ.get('KRATOS_DSN', ''),
    'KRATOS_CIPHER_SECRET': os.environ.get('KRATOS_CIPHER_SECRET', ''),
    'KRATOS_COOKIE_SECRET': os.environ.get('KRATOS_COOKIE_SECRET', ''),
    'KRATOS_SMTP_CONNECTION_URI': os.environ.get('KRATOS_SMTP_CONNECTION_URI', ''),
}

with open(template_file, 'r') as f:
    content = f.read()

for var_name, var_value in vars_to_sub.items():
    placeholder = '${' + var_name + '}'
    content = content.replace(placeholder, var_value)

with open(output_file, 'w') as f:
    f.write(content)

print(f'Substituted variables in {template_file} -> {output_file}')
