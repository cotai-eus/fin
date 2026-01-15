/**
 * src/shared/components/PageHeader.tsx
 *
 * Page Header Container Component
 *
 * Standard page header for consistent styling across dashboard pages
 */

interface PageHeaderProps {
  title: string;
  subtitle?: string;
  description?: string;
  action?: React.ReactNode;
}

export function PageHeader({
  title,
  subtitle,
  description,
  action,
}: PageHeaderProps) {
  return (
    <div className="flex items-start justify-between">
      <div className="flex-1">
        <h1 className="text-3xl font-bold text-gray-900">{title}</h1>
        {subtitle && <p className="mt-1 text-lg text-gray-600">{subtitle}</p>}
        {description && <p className="mt-2 text-gray-600">{description}</p>}
      </div>
      {action && <div className="ml-4">{action}</div>}
    </div>
  );
}
