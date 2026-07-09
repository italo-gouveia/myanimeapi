import type { Character } from '../lib/api/types'

interface CharacterCardProps {
  character: Character
}

function CharacterAvatar({ character }: { character: Character }) {
  if (character.image_url) {
    return (
      <img
        src={character.image_url}
        alt={character.name}
        className="w-full h-full object-cover"
        loading="lazy"
      />
    )
  }
  const initials = character.name
    .split(' ')
    .slice(0, 2)
    .map((w) => w[0]?.toUpperCase() ?? '')
    .join('')
  return (
    <div className="w-full h-full bg-gradient-to-br from-brand-100 to-brand-200 flex items-center justify-center">
      <span className="text-2xl font-bold text-brand-600 select-none">{initials}</span>
    </div>
  )
}

export function CharacterCard({ character }: CharacterCardProps) {
  return (
    <div
      className="flex flex-col bg-white border border-slate-200 rounded-lg overflow-hidden hover:shadow-md hover:border-brand-300 transition-shadow"
      data-testid={`character-card-${character.id}`}
    >
      <div className="relative w-full aspect-[3/4] bg-slate-100 overflow-hidden">
        <CharacterAvatar character={character} />
      </div>

      <div className="p-3 flex flex-col gap-1">
        <h3 className="font-semibold text-slate-900 line-clamp-2 text-sm leading-snug">
          {character.name}
        </h3>
        {character.voice_actor && (
          <p className="text-xs text-slate-500 truncate" title={character.voice_actor}>
            🎙 {character.voice_actor}
          </p>
        )}
        {character.description && (
          <p className="text-xs text-slate-400 line-clamp-2 mt-0.5">{character.description}</p>
        )}
      </div>
    </div>
  )
}
