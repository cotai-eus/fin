#!/bin/sh
set -e

# Escape special sed characters in replacement strings
escape_sed() {
  printf '%s\n' "$1" | sed -e 's/[\/&]/\\&/g'
}

# Read the template and perform substitutions
{
  while IFS= read -r line; do
    line="${line%$'\n'}"
    # Replace each placeholder
    line=$(printf '%s\n' "$line" | sed "s/\${KRATOS_DSN}/$(escape_sed "$KRATOS_DSN")/g")
    line=$(printf '%s\n' "$line" | sed "s/\${KRATOS_CIPHER_SECRET}/$(escape_sed "$KRATOS_CIPHER_SECRET")/g")
    line=$(printf '%s\n' "$line" | sed "s/\${KRATOS_COOKIE_SECRET}/$(escape_sed "$KRATOS_COOKIE_SECRET")/g")
    line=$(printf '%s\n' "$line" | sed "s/\${KRATOS_SMTP_CONNECTION_URI}/$(escape_sed "$KRATOS_SMTP_CONNECTION_URI")/g")
    echo "$line"
  done
} < /etc/config/kratos/kratos.yml.template > /etc/config/kratos/kratos.yml

# Run kratos with all arguments
exec kratos "$@"



