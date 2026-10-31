import { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import {
  listCharacters,
  createCharacter,
  updateCharacter,
  deleteCharacter,
  type CreateCharacterPayload,
} from '../../lib/api/characters'
import type { Character } from '../../lib/api/types'

const LIMIT = 20

interface CharacterFormProps {
  initial?: Character
  onSubmit: (payload: CreateCharacterPayload) => void
  onCancel: () => void
  isPending: boolean
}

function CharacterForm({ initial, onSubmit, onCancel, isPending }: CharacterFormProps) {
  const [name, setName] = useState(initial?.name ?? '')
  const [description, setDescription] = useState(initial?.description ?? '')
  const [voiceActor, setVoiceActor] = useState(initial?.voice_actor ?? '')
  const [imageUrl, setImageUrl] = useState(initial?.image_url ?? '')

  function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    onSubmit({ name, description, voice_actor: voiceActor, image_url: imageUrl })
  }

  return (
    <form onSubmit={handleSubmit} className="flex flex-col gap-4">
      <div className="flex flex-col gap-1">
        <label className="text-sm font-medium text-slate-700" htmlFor="char-name">
          Name <span className="text-red-500">*</span>
        </label>
        <input
          id="char-name"
          required
          value={name}
          onChange={(e) => setName(e.target.value)}
          className="border border-slate-200 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-brand-500"
          placeholder="Character name"
        />
      </div>
      <div className="flex flex-col gap-1">
        <label className="text-sm font-medium text-slate-700" htmlFor="char-desc">
          Description
        </label>
        <textarea
          id="char-desc"
          rows={3}
          value={description}
          onChange={(e) => setDescription(e.target.value)}
          className="border border-slate-200 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-brand-500 resize-none"
          placeholder="Short description"
        />
      </div>
      <div className="flex flex-col gap-1">
        <label className="text-sm font-medium text-slate-700" htmlFor="char-va">
          Voice Actor
        </label>
        <input
          id="char-va"
          value={voiceActor}
          onChange={(e) => setVoiceActor(e.target.value)}
          className="border border-slate-200 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-brand-500"
          placeholder="Voice actor name"
        />
      </div>
      <div className="flex flex-col gap-1">
        <label className="text-sm font-medium text-slate-700" htmlFor="char-img">
          Image URL
        </label>
        <input
          id="char-img"
          type="url"
          value={imageUrl}
          onChange={(e) => setImageUrl(e.target.value)}
          className="border border-slate-200 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-brand-500"
          placeholder="https://..."
        />
      </div>
      <div className="flex justify-end gap-3 pt-2">
        <button
          type="button"
          onClick={onCancel}
          className="px-4 py-2 text-sm rounded-lg border border-slate-200 hover:bg-slate-50"
        >
          Cancel
        </button>
        <button
          type="submit"
          disabled={isPending}
          className="px-4 py-2 text-sm rounded-lg bg-brand-600 text-white hover:bg-brand-700 disabled:opacity-50"
        >
          {isPending ? 'Saving…' : initial ? 'Save changes' : 'Create'}
        </button>
      </div>
    </form>
  )
}

interface ModalProps {
  title: string
  onClose: () => void
  children: React.ReactNode
}

function Modal({ title, onClose, children }: ModalProps) {
  return (
    <div
      className="fixed inset-0 z-50 flex items-center justify-center bg-black/40"
      role="dialog"
      aria-modal="true"
      aria-labelledby="modal-title"
    >
      <div className="bg-white rounded-xl shadow-xl w-full max-w-md mx-4 p-6 flex flex-col gap-4">
        <div className="flex items-center justify-between">
          <h2 id="modal-title" className="text-lg font-semibold text-slate-800">
            {title}
          </h2>
          <button
            type="button"
            onClick={onClose}
            className="text-slate-400 hover:text-slate-600 text-xl leading-none"
            aria-label="Close"
          >
            ×
          </button>
        </div>
        {children}
      </div>
    </div>
  )
}

type ModalState =
  | { kind: 'closed' }
  | { kind: 'create' }
  | { kind: 'edit'; character: Character }
  | { kind: 'delete'; character: Character }

export function AdminCharactersPage() {
  const [page, setPage] = useState(1)
  const [modal, setModal] = useState<ModalState>({ kind: 'closed' })
  const queryClient = useQueryClient()

  const { data, isLoading, isError } = useQuery({
    queryKey: ['admin', 'characters', page],
    queryFn: () => listCharacters({ page, limit: LIMIT }),
  })

  const createMutation = useMutation({
    mutationFn: (payload: CreateCharacterPayload) => createCharacter(payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['admin', 'characters'] })
      setModal({ kind: 'closed' })
    },
  })

  const updateMutation = useMutation({
    mutationFn: ({ id, payload }: { id: number; payload: Partial<CreateCharacterPayload> }) =>
      updateCharacter(id, payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['admin', 'characters'] })
      setModal({ kind: 'closed' })
    },
  })

  const deleteMutation = useMutation({
    mutationFn: (id: number) => deleteCharacter(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['admin', 'characters'] })
      setModal({ kind: 'closed' })
    },
  })

  const totalPages = data ? Math.ceil(data.total / LIMIT) : 1

  return (
    <div className="flex flex-col gap-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold">Characters</h1>
          {data && (
            <span className="text-sm text-slate-500">{data.total.toLocaleString()} total</span>
          )}
        </div>
        <button
          type="button"
          onClick={() => setModal({ kind: 'create' })}
          className="px-4 py-2 text-sm rounded-lg bg-brand-600 text-white hover:bg-brand-700"
        >
          + New character
        </button>
      </div>

      {isError && (
        <p className="text-red-600" role="alert">
          Failed to load characters.
        </p>
      )}

      {(createMutation.isError || updateMutation.isError || deleteMutation.isError) && (
        <p className="text-red-600" role="alert">
          Operation failed. Please try again.
        </p>
      )}

      <div className="bg-white border border-slate-200 rounded-xl overflow-hidden">
        <table className="w-full text-sm">
          <thead className="bg-slate-50 border-b border-slate-200">
            <tr>
              <th className="text-left px-4 py-3 font-semibold text-slate-600 w-12">ID</th>
              <th className="text-left px-4 py-3 font-semibold text-slate-600">Name</th>
              <th className="text-left px-4 py-3 font-semibold text-slate-600 hidden md:table-cell">
                Voice Actor
              </th>
              <th className="text-left px-4 py-3 font-semibold text-slate-600 hidden lg:table-cell">
                Description
              </th>
              <th className="text-right px-4 py-3 font-semibold text-slate-600">Actions</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-slate-100">
            {isLoading
              ? Array.from({ length: 8 }).map((_, i) => (
                  <tr key={i}>
                    {Array.from({ length: 5 }).map((__, j) => (
                      <td key={j} className="px-4 py-3">
                        <div className="h-4 bg-slate-100 rounded animate-pulse" />
                      </td>
                    ))}
                  </tr>
                ))
              : data?.data.map((char) => (
                  <tr key={char.id} className="hover:bg-slate-50 transition-colors">
                    <td className="px-4 py-3 text-slate-400 tabular-nums">{char.id}</td>
                    <td className="px-4 py-3 font-medium text-slate-800">
                      <div className="flex items-center gap-2">
                        {char.image_url && (
                          <img
                            src={char.image_url}
                            alt=""
                            className="w-7 h-7 rounded-full object-cover shrink-0"
                          />
                        )}
                        {char.name}
                      </div>
                    </td>
                    <td className="px-4 py-3 text-slate-500 hidden md:table-cell">
                      {char.voice_actor ?? '—'}
                    </td>
                    <td className="px-4 py-3 text-slate-500 hidden lg:table-cell max-w-xs truncate">
                      {char.description ?? '—'}
                    </td>
                    <td className="px-4 py-3 text-right">
                      <div className="flex items-center justify-end gap-2">
                        <button
                          type="button"
                          onClick={() => setModal({ kind: 'edit', character: char })}
                          className="px-3 py-1 text-xs rounded border border-slate-200 hover:bg-slate-50"
                        >
                          Edit
                        </button>
                        <button
                          type="button"
                          onClick={() => setModal({ kind: 'delete', character: char })}
                          className="px-3 py-1 text-xs rounded border border-red-200 text-red-600 hover:bg-red-50"
                        >
                          Delete
                        </button>
                      </div>
                    </td>
                  </tr>
                ))}
          </tbody>
        </table>
      </div>

      {totalPages > 1 && (
        <div className="flex items-center justify-between">
          <button
            type="button"
            disabled={page === 1}
            onClick={() => setPage((p) => p - 1)}
            className="px-3 py-1.5 rounded border border-slate-200 text-sm disabled:opacity-40 hover:bg-slate-50"
          >
            ← Previous
          </button>
          <span className="text-sm text-slate-500">
            Page {page} of {totalPages}
          </span>
          <button
            type="button"
            disabled={page === totalPages}
            onClick={() => setPage((p) => p + 1)}
            className="px-3 py-1.5 rounded border border-slate-200 text-sm disabled:opacity-40 hover:bg-slate-50"
          >
            Next →
          </button>
        </div>
      )}

      {modal.kind === 'create' && (
        <Modal title="New character" onClose={() => setModal({ kind: 'closed' })}>
          <CharacterForm
            onSubmit={(payload) => createMutation.mutate(payload)}
            onCancel={() => setModal({ kind: 'closed' })}
            isPending={createMutation.isPending}
          />
        </Modal>
      )}

      {modal.kind === 'edit' && (
        <Modal title="Edit character" onClose={() => setModal({ kind: 'closed' })}>
          <CharacterForm
            initial={modal.character}
            onSubmit={(payload) => updateMutation.mutate({ id: modal.character.id, payload })}
            onCancel={() => setModal({ kind: 'closed' })}
            isPending={updateMutation.isPending}
          />
        </Modal>
      )}

      {modal.kind === 'delete' && (
        <Modal title="Delete character" onClose={() => setModal({ kind: 'closed' })}>
          <p className="text-slate-600 text-sm">
            Delete <strong>{modal.character.name}</strong>? This cannot be undone.
          </p>
          <div className="flex justify-end gap-3">
            <button
              type="button"
              onClick={() => setModal({ kind: 'closed' })}
              className="px-4 py-2 text-sm rounded-lg border border-slate-200 hover:bg-slate-50"
            >
              Cancel
            </button>
            <button
              type="button"
              disabled={deleteMutation.isPending}
              onClick={() => deleteMutation.mutate(modal.character.id)}
              className="px-4 py-2 text-sm rounded-lg bg-red-600 text-white hover:bg-red-700 disabled:opacity-50"
            >
              {deleteMutation.isPending ? 'Deleting…' : 'Delete'}
            </button>
          </div>
        </Modal>
      )}
    </div>
  )
}
