import { component$ } from "@builder.io/qwik";
import { Link } from "@builder.io/qwik-city";

type BreadcrumbsProps = {
  crumbs: { text: string; path?: string }[];
};

export const Breadcrumbs = component$<BreadcrumbsProps>(({ crumbs }) => (
  <div class="breadcrumbs text-sm">
    <ul>
      {crumbs.map((crumb) => (
        <li key={crumb.text}>
          {crumb.path ? (
            <Link href={crumb.path}>{crumb.text}</Link>
          ) : (
            <>{crumb.text}</>
          )}
        </li>
      ))}
    </ul>
  </div>
));
