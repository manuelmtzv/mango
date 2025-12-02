import { noteRoutes } from "@/modules/note/routes/note.routes";
import { tagRoutes } from "@/modules/tags/routes/tag.routes";

export const appRoutes = [...noteRoutes, ...tagRoutes];
