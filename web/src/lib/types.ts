export type Role = 'customer' | 'organizer' | 'admin';
export interface Account {
  id: string;
  role: Role;
  name: string;
  email: string;
  city: string;
  bio: string;
  avatar_url: string;
  interests: string[];
  preferences: Record<string, boolean>;
  onboarding_done: boolean;
  organizer_id?: string;
  organizer_slug?: string;
}
export interface Organizer {
  id: string;
  slug: string;
  name: string;
  description: string;
  about: string;
  city: string;
  category: string;
  avatar_url: string;
  cover_url: string;
  followers: number;
  verified: boolean;
}
export interface Mention {
  kind: 'event' | 'organizer';
  id: string;
  slug: string;
  name: string;
  city?: string;
  label?: string;
}
export interface EventSession {
  id: string;
  event_id?: string;
  label: string;
  starts_at: string;
  ends_at: string;
  price: number;
  capacity: number;
  available: number;
}
export interface ArtsEvent {
  id: string;
  slug: string;
  title: string;
  description: string;
  category: string;
  city: string;
  organizer_id: string;
  organizer: Organizer;
  venue: string;
  address: string;
  maps_url: string;
  duration: string;
  language: string;
  age: string;
  flyer_url: string;
  trailer_url: string;
  layout_url: string;
  lineup: { name: string; role: string; photo: string }[];
  published: boolean;
  starts_at: string;
  ends_at: string;
  price: number;
  sessions: EventSession[];
}
export interface Post {
  id: string;
  account_id: string;
  organizer_id?: string;
  author: string;
  avatar_url: string;
  official: boolean;
  title: string;
  body: string;
  image_url: string;
  link_url: string;
  link_label: string;
  city: string;
  pinned: boolean;
  created_at: string;
  votes: number;
  comment_count: number;
  mentions: Mention[];
}
export interface Review {
  id: string;
  event_id: string;
  event_title: string;
  event_slug: string;
  event_starts_at?: string;
  organizer_name?: string;
  organizer_id: string;
  account_id: string;
  author: string;
  avatar_url: string;
  body: string;
  created_at: string;
  updated_at: string;
  verified: boolean;
  helpful: number;
  replies: { id: string; author: string; body: string; created_at: string }[];
}
export interface PaymentMethod {
  id: string;
  organizer_id: string;
  kind: 'bank' | 'wallet';
  provider: string;
  number: string;
  holder: string;
  note: string;
  enabled: boolean;
}
export type OrderStatus =
  | 'awaiting_payment'
  | 'awaiting_review'
  | 'correction_requested'
  | 'approved'
  | 'cancelled'
  | 'expired';
export interface Ticket {
  id: string;
  order_id: string;
  ordinal: number;
  qr: string;
  checked_at?: string;
  event?: ArtsEvent;
  session?: EventSession;
}
export interface Order {
  id: string;
  account_id: string;
  session_id: string;
  quantity: number;
  total: number;
  status: OrderStatus;
  expires_at?: string;
  payment_snapshot?: PaymentMethod;
  proof_id?: string;
  note: string;
  event: ArtsEvent;
  session: EventSession;
  customer: { name: string; email: string };
  tickets: Ticket[];
  history: { id: number; status: OrderStatus; note: string; created_at: string }[];
  created_at: string;
}
export interface SessionState {
  user: Account | null;
  follows: string[];
  votes: string[];
  helpful_reviews: string[];
  bookmarks: { kind: string; id: string }[];
  unread: number;
}
export interface AppContext {
  readonly account: SessionState;
  readonly user: Account | null;
  readonly city: string;
  readonly demo: boolean;
  toast: (message: string) => void;
  requireLogin: () => boolean;
  refresh: () => Promise<void>;
}
