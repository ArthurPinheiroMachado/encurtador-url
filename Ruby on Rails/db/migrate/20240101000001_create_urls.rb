class CreateUrls < ActiveRecord::Migration[7.1]
  def change
    create_table :url, id: false do |t|
      t.string :id, primary_key: true, limit: 8, null: false
      t.text :original, null: false
      t.bigint :accesses, default: 0
    end

    add_index :url, :original, unique: true
  end
end
