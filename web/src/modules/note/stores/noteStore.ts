import { ref, computed } from "vue";
import { defineStore } from "pinia";
import { INote } from "@/modules/note/interfaces/INote";
import noteDefaults from "@/modules/note/config/defaults";

export const useNoteStore = defineStore("noteStore", () => {
  const notes = ref<INote[]>([]);

  return {
    notes,

    sortedByTime: computed(() => {
      return notes.value.slice().sort((a, b) => {
        return (
          new Date(b.updatedAt).getTime() - new Date(a.updatedAt).getTime()
        );
      });
    }),

    getNotes(): void {},
    addNote(note: INote): void {
      const { NOTE_LIMIT } = noteDefaults;

      if (notes.value.length >= NOTE_LIMIT) {
        notes.value = notes.value.slice(0, NOTE_LIMIT - 1);
      }

      const updatedNotes = [...notes.value, note];
      notes.value = updatedNotes;
    },
    setNotes(newNotes: INote[] | undefined): void {
      if (newNotes) notes.value = newNotes;
    },
    getNoteById(id: string): INote {
      const note = notes.value.filter((note) => {
        return note.id === id;
      })[0];

      return note || null;
    },
    updateNote(id: string, title: string, content: string): void {
      const note = this.getNoteById(id);

      if (note) {
        note.title = title;
        note.content = content;
        note.updatedAt = new Date().toISOString();
      }
    },
    deleteNote(id: string): void {
      const note = this.getNoteById(id);

      notes.value = notes.value.filter((storeNote) => {
        return storeNote.id != note.id;
      });
    },
    resetState(): void {
      notes.value = [];
    },
  };
});
